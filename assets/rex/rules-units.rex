^
[^:]+:      # filename, e.g. "vehicles.yaml:"
[^.]+[.]    # unit name, e.g. "IFV."
(
	Tooltip(@[^.]+)?[.]Name     # "Tooltip(@acolturr).Name"
	|Buildable[.]Description    # "Buildable.Description"
	|TooltipDescription(@[^.]+)?[.]Description  # "TooltipDescription(@ally).Description"
	|TooltipExtras(@[^.]+)?[.][^.]+             # "TooltipExtras.(Attributes|Description|Strengths|Weaknesses)"
)
$
