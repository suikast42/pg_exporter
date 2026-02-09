package exporter

import (
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

// normalizeKingpinBoolEqualsArgs rewrites boolean flags passed as `--flag=false`
// (or `-f=false`) into kingpin-compatible forms (`--no-flag`).
//
// kingpin bool flags are "presence flags" and don't accept values, which leads
// to confusing parse errors like: "unexpected false, try --help".
//
// This function is intentionally conservative: it only rewrites flags that are
// known boolean flags in the provided kingpin model.
func normalizeKingpinBoolEqualsArgs(args []string, model *kingpin.ApplicationModel) []string {
	if len(args) == 0 || model == nil || model.FlagGroupModel == nil {
		return args
	}

	boolFlags := make(map[string]struct{})
	shortToLong := make(map[byte]string)
	for _, f := range model.FlagGroupModel.Flags {
		if f == nil || !f.IsBoolFlag() {
			continue
		}
		boolFlags[f.Name] = struct{}{}
		if f.Short != 0 && f.Short <= 0x7f { // only ASCII shorts are relevant to CLI parsing
			shortToLong[byte(f.Short)] = f.Name
		}
	}
	if len(boolFlags) == 0 {
		return args
	}

	out := make([]string, 0, len(args))
	for _, arg := range args {
		// Long form: --flag=false
		if strings.HasPrefix(arg, "--") {
			name, val, ok := strings.Cut(arg[2:], "=")
			if ok {
				if _, isBool := boolFlags[name]; isBool {
					if b, err := strconv.ParseBool(val); err == nil {
						if b {
							out = append(out, "--"+name)
						} else {
							out = append(out, "--no-"+name)
						}
						continue
					}
				}
			}
			out = append(out, arg)
			continue
		}

		// Short form: -f=false (only single short flag)
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
			name, val, ok := strings.Cut(arg[1:], "=")
			if ok && len(name) == 1 {
				if long, exists := shortToLong[name[0]]; exists {
					if _, isBool := boolFlags[long]; isBool {
						if b, err := strconv.ParseBool(val); err == nil {
							if b {
								out = append(out, "-"+name)
							} else {
								out = append(out, "--no-"+long)
							}
							continue
						}
					}
				}
			}
			out = append(out, arg)
			continue
		}

		out = append(out, arg)
	}
	return out
}

