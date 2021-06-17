package dwdparse

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"stash.us.cray.com/dpm/dws-operator/pkg/apis/dws/v1alpha1"
	client "stash.us.cray.com/dpm/dws-operator/pkg/client/clientset/versioned/typed/dws/v1alpha1"
)

type dwUnsupportedCommandErr struct {
	command string
}

func NewUnsupportedCommandErr(command string) error {
	return &dwUnsupportedCommandErr{command}
}

func (e *dwUnsupportedCommandErr) Error() string {
	return fmt.Sprintf("Unsupported Command: '%s'", e.command)
}

// IsUnsupportedCommand returns true if the error indicates that the command
// is unsupported
func IsUnsupportedCommand(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*dwUnsupportedCommandErr)
	return ok
}

// BuildRulesMap builds a map of the DWDirectives argument parser rules for the specified command
func BuildRulesMap(rules []v1alpha1.DWDirectiveRuleSpec, cmd string) (map[string]v1alpha1.DWDirectiveRuleDef, error) {
	rulesMap := make(map[string]v1alpha1.DWDirectiveRuleDef)

	// Search for the command in the supported commands within the ruleset
	for _, r := range rules {
		if cmd == r.Command {
			for _, rd := range r.RuleDefs {
				rulesMap[rd.Key] = rd
			}
		}
	}

	if len(rulesMap) == 0 {
		return nil, NewUnsupportedCommandErr(cmd)
	}

	return rulesMap, nil
}

// BuildArgsMap builds a map of the DWDirective's arguments in the form: args["key"] = value
func BuildArgsMap(dwd string) (map[string]string, error) {
	argsMap := make(map[string]string)
	dwdArgs := strings.Fields(dwd)
	if dwdArgs[0] == "#DW" {
		argsMap["command"] = dwdArgs[1]
		for i := 2; i < len(dwdArgs); i++ {
			keyValue := strings.Split(dwdArgs[i], "=")

			// Don't allow repeated arguments
			_, ok := argsMap[keyValue[0]]
			if ok {
				return nil, errors.New("repeated argument in directive: " + keyValue[0])
			}

			if len(keyValue) == 1 {
				argsMap[keyValue[0]] = "true"
			} else if len(keyValue) == 2 {
				argsMap[keyValue[0]] = keyValue[1]
			} else {
				keyValue := strings.SplitN(dwdArgs[i], "=", 2)
				argsMap[keyValue[0]] = keyValue[1]
			}
		}
	} else {
		return nil, errors.New("missing #DW in directive")
	}
	return argsMap, nil
}

// ValidateArgs validates a map of arguments against the rules
// For cases where an unknown command may be allowed because there may be other handlers for that command
//   failUnknownCommand = false
func ValidateArgs(args map[string]string, rules []v1alpha1.DWDirectiveRuleSpec, failUnknownCommand bool) error {
	command := args["command"]

	// Determine the rules map for command
	rulesMap, err := BuildRulesMap(rules, command)
	if err != nil {
		// If the command is unsupported and we are supposed to fail in that case return error.
		// Otherwise just return nil to effectively skip the #DW
		// for info on errors.As() below see:
		// https://stackoverflow.com/questions/62441960/error-wrap-unwrap-type-checking-with-errors-is#62442136
		var unsupportedCommand *dwUnsupportedCommandErr
		if failUnknownCommand && errors.As(err, &unsupportedCommand) {
			return err
		} else {
			return nil
		}
	}

	// Compile this regex outside the loop for better performance.
	var boolMatcher = regexp.MustCompile(`(?i)^(true|false)$`) // (?i) -> case-insensitve comparison

	// Create a map that maps a directive rule definition to an argument that correctly matches it
	// key: DWDirectiveRule	value: argument that matches that rule
	// Required to check that all DWDirectiveRuleDef's have been met
	argToRuleMap := map[v1alpha1.DWDirectiveRuleDef]string{}

	// Iterate over all arguments and validate each based on the associated rule
	for k, v := range args {
		if k != "command" {
			rule, found := rulesMap[k]
			if !found {
				return errors.New("unsupported argument - " + k)
			}
			if rule.IsValueRequired && len(v) == 0 {
				return errors.New("malformed keyword[=value]: " + k + "=" + v)
			}
			switch rule.Type {
			case "integer":
				// i,err := strconv.ParseInt(v, 10, 64)
				i, err := strconv.Atoi(v)
				if err != nil {
					return errors.New("invalid integer argument: " + k + "=" + v)
				}
				if rule.Max != 0 && i > rule.Max {
					return errors.New("specified integer exceeds maximum " + strconv.Itoa(rule.Max) + ": " + k + "=" + v)
				}
				if rule.Min != 0 && i < rule.Min {
					return errors.New("specified integer smaller than minimum " + strconv.Itoa(rule.Min) + ": " + k + "=" + v)
				}
			case "bool":
				if rule.Pattern != "" {
					isok := boolMatcher.MatchString(v)
					if !isok {
						return errors.New("invalid bool argument: " + k + "=" + v)
					}
				}
			case "string":
				if rule.Pattern != "" {
					isok, err := regexp.MatchString(rule.Pattern, v)
					if !isok {
						if err != nil {
							return errors.New("invalid regexp in rule: " + rule.Pattern)
						}
						return errors.New("invalid argument: " + k + "=" + v)
					}
				}
			default:
				return errors.New("unsupported value type: " + rule.Type)
			}

			// NOTE: We know that we don't have repeated arguments here because the arguments
			//       come to us in a map indexed by the argment name.
			argToRuleMap[rule] = k
		}
	}

	// Iterate over the rules to ensure all required rules have an argument
	for k, v := range rulesMap {
		// Ensure that each required rule has an argument
		if v.IsRequired {
			_, ok := argToRuleMap[v]
			if !ok {
				return errors.New("missing argument: " + k)
			}
		}
	}

	return nil
}

// GetParserRules is used to get the DWDirective parser rule set
func GetParserRules(ruleSetName string, namespace string) (*v1alpha1.DWDirectiveRule, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	dwsClient, err := client.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dwdRules, err := dwsClient.DWDirectiveRules(namespace).Get(ruleSetName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return dwdRules, nil
}

// ValidateDWDirectives validates a set of #DW directives against a specified rule set
func ValidateDWDirectives(directives []string, ruleSetName string, namespace string, failUnknownCommand bool) error {

	dwdRules, err := GetParserRules(ruleSetName, namespace)
	if err != nil {
		return err
	}

	for _, dwd := range directives {
		// Build a map of the #DW commands and arguments
		argsMap, err := BuildArgsMap(dwd)
		if err != nil {
			return err
		}

		err = ValidateArgs(argsMap, dwdRules.Spec, failUnknownCommand)
		if err != nil {
			return err
		}
	}

	return nil
}
