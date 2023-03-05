package cmd

import (
	"fmt"
	"io"
	// "os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/cobra"
	"github.com/zclconf/go-cty/cty"
)

type Template struct {
	Destination string `hcl:"destination"`
	Contents    string `hcl:"contents"`
	Source      string `hcl:"source"`
}

func init() {
	rootCmd.AddCommand(templateBlockCmd())
}

func templateBlockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "templateblock",
		Short: "Edit template block from stdin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		templateBlockUpsertCmd(),
	)

	return cmd
}

func templateBlockUpsertCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "upsert",
		Short: "add/edit template block from stdin to stdout",
		Long: `
Search for a template block in stdin with the same 'destinaton' attribute and edit it.
If no template block with same 'destination' attribute is found, a new one is created.

Info:
if you edit a template block with e.g. a destination and contents attribute, the source attribute will be removed, and vice versa.
This is because the vault agent only supports one of the two attributes.
https://developer.hashicorp.com/vault/docs/agent/template#templating-configuration-example
`,
		RunE: runTemplateBlockUpsertCmd,
	}
	// Add flags only for this command
	cmd.Flags().StringP("destination", "d", "", "destination of the template block (Required)")
	cmd.Flags().StringP("source", "s", "", "source of the template block (Required if contents is not provided)")
	cmd.Flags().StringP("contents", "c", "", "contents of the template block (Required if source is not provided)")

	return cmd
}

func runTemplateBlockUpsertCmd(cmd *cobra.Command, args []string) error {
	destination, err := cmd.Flags().GetString("destination")
	if err != nil {
		return err
	}
	if destination == "" {
		return fmt.Errorf("destination is required")
	}

	source, err := cmd.Flags().GetString("source")
	if err != nil {
		return err
	}
	contents, err := cmd.Flags().GetString("contents")
	if err != nil {
		return err
	}
	if source == "" && contents == "" {
		return fmt.Errorf("either source or contents is required")
	}

	var input []byte
	input, err = io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return err
	}

	// Parse stdin as hcl file
	parser := hclparse.NewParser()
	parsedInput, parseDiags := parser.ParseHCL(input, "vault-agent.hcl")
	if parseDiags.HasErrors() {
		return fmt.Errorf(parseDiags.Error())
	}

	// Parse hcl file as hclwrite file
	config, diags := hclwrite.ParseConfig(parsedInput.Bytes, "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return fmt.Errorf(diags.Error())
	}

	// check if a template block with the same destination exists and replace it with the new one
	found := false
	for _, block := range config.Body().Blocks() {
		//Check if block is a template block
		if block.Type() != "template" {
			continue
		}
		//Check if block has a destination attribute
		destinationAttribut := block.Body().GetAttribute("destination")
		if destinationAttribut == nil {
			continue
		}

		//Remove block if destination is the same as the new template Block
		destinationString := string(destinationAttribut.Expr().BuildTokens(nil).Bytes())

		//Remove Whitescpaces and Quotes at the beginning and end of the string
		destinationString = strings.TrimSpace(destinationString)
		destinationString = strings.TrimRight(strings.TrimLeft(destinationString, "\""), "\"")

		if destinationString == destination {
			block.Body().SetAttributeValue("destination", cty.StringVal(destination))
			block.Body().RemoveAttribute("contents")
			block.Body().RemoveAttribute("source")
			if contents != "" {
				block.Body().SetAttributeValue("contents", cty.StringVal(contents))
			}
			if source != "" {
				block.Body().SetAttributeValue("source", cty.StringVal(source))
			}
			found = true
			break
		}
	}

	// Add new template block when no template block with the same destination was found
	if !found {
		if len(config.Body().Blocks()) > 0 {
			config.Body().AppendNewline()
		}
		templateBlock := config.Body().AppendNewBlock("template", nil)
		templateBlock.Body().SetAttributeValue("destination", cty.StringVal(destination))
		if contents != "" {
			templateBlock.Body().SetAttributeValue("contents", cty.StringVal(contents))
		}
		if source != "" {
			templateBlock.Body().SetAttributeValue("source", cty.StringVal(source))
		}
	}

	//Format Output
	raw := config.BuildTokens(nil).Bytes()
	out := hclwrite.Format(raw)
	//Write to out to stdout of cmd
	cmd.OutOrStdout().Write(out)

	// os.Stdout.Write(out)
	return nil
}

func getTemplateBlockUpsertCmdFlags(cmd *cobra.Command) (Template, error) {
	var template Template
	destination, err := cmd.Flags().GetString("destination")
	if err != nil {
		return template, err
	}
	if destination == "" {
		return template, fmt.Errorf("destination is required")
	}

	source, err := cmd.Flags().GetString("source")
	if err != nil {
		return template, err
	}
	contents, err := cmd.Flags().GetString("contents")
	if err != nil {
		return template, err
	}
	if source == "" && contents == "" {
		return template, fmt.Errorf("either source or contents is required")
	}

	template.Destination = destination
	template.Contents = contents
	template.Source = source

	return template, nil

}
