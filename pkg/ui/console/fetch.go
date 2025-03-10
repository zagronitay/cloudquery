package console

import (
	"sort"

	"github.com/cloudquery/cloudquery/pkg/client"
	"github.com/cloudquery/cloudquery/pkg/ui"
	"github.com/cloudquery/cq-provider-sdk/provider/schema/diag"
)

func printFetchResponse(summary *client.FetchResponse) {
	if summary == nil {
		return
	}
	for _, pfs := range summary.ProviderFetchSummary {
		if len(pfs.Diagnostics()) > 0 {
			printDiagnostics(pfs.ProviderName, pfs.Diagnostics())
			continue
		}
		if len(pfs.PartialFetchErrors) == 0 {
			continue
		}
		ui.ColorizedOutput(ui.ColorHeader, "Partial Fetch Errors for Provider %s:\n\n", pfs.ProviderName)
		for _, r := range pfs.PartialFetchErrors {
			if r.RootTableName != "" {
				ui.ColorizedOutput(ui.ColorErrorBold,
					"Parent-Resource: %-64s Parent-Primary-Keys: %v, Table: %s, Error: %s\n",
					r.RootTableName,
					r.RootPrimaryKeyValues,
					r.TableName,
					r.Error)
			} else {
				ui.ColorizedOutput(ui.ColorErrorBold,
					"Table: %-64s Error: %s\n",
					r.TableName,
					r.Error)
			}
		}
		ui.ColorizedOutput(ui.ColorWarning, "\n")
	}
}

func printDiagnostics(providerName string, diags diag.Diagnostics) {
	// sort diagnostics by severity/type
	sort.Sort(diags)
	ui.ColorizedOutput(ui.ColorHeader, "Fetch Diagnostics for provider %s:\n\n", providerName)
	for _, d := range diags {
		desc := d.Description()
		switch d.Severity() {
		case diag.IGNORE:
			ui.ColorizedOutput(ui.ColorHeader, "Resource: %-10s Type: %-10s Severity: %s\n\tSummary: %s\n",
				ui.ColorProgress.Sprintf("%s", desc.Resource),
				ui.ColorProgressBold.Sprintf("%s", d.Type()),
				ui.ColorDebug.Sprintf("Ignore"),
				ui.ColorDebug.Sprintf("%s", desc.Summary))
		case diag.WARNING:
			ui.ColorizedOutput(ui.ColorHeader, "Resource: %-10s Type: %-10s Severity: %s\n\tSummary: %s\n",
				ui.ColorInfo.Sprintf("%s", desc.Resource),
				ui.ColorProgressBold.Sprintf("%s", d.Type()),
				ui.ColorWarning.Sprintf("Warning"),
				ui.ColorWarning.Sprintf("%s", desc.Summary))
		case diag.ERROR:
			ui.ColorizedOutput(ui.ColorHeader, "Resource: %-10s Type: %-10s Severity: %s\n\tSummary: %s\n",
				ui.ColorProgress.Sprintf("%s", desc.Resource),
				ui.ColorProgressBold.Sprintf("%s", d.Type()),
				ui.ColorErrorBold.Sprintf("Error"),
				ui.ColorErrorBold.Sprintf("%s", desc.Summary))
		}
		if desc.Detail != "" {
			ui.ColorizedOutput(ui.ColorInfo, "\tRemediation: %s\n", desc.Detail)
		}
	}
	ui.ColorizedOutput(ui.ColorInfo, "\n")
}
