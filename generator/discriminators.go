package generator

import (
	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
)

type discInfo struct {
	Discriminators map[string]discor
	Discriminated  map[string]discee
}

type discor struct {
	FieldName string   `json:"fieldName"`
	GoType    string   `json:"goType"`
	JSONName  string   `json:"jsonName"`
	Children  []discee `json:"children"`
}

type discee struct {
	FieldName  string   `json:"fieldName"`
	FieldValue string   `json:"fieldValue"`
	GoType     string   `json:"goType"`
	JSONName   string   `json:"jsonName"`
	Ref        spec.Ref `json:"ref"`
	ParentRef  spec.Ref `json:"parentRef"`
}

func discriminatorInfo(doc *spec.Document) *discInfo {
	baseTypes := make(map[string]discor)
	for _, sch := range doc.AllDefinitions() {
		if sch.Schema.Discriminator != "" {
			baseTypes[sch.Ref.String()] = discor{
				// TODO: more trickery to allow for name customization
				FieldName: sch.Schema.Discriminator,
				GoType:    swag.ToGoName(sch.Name),
				JSONName:  sch.Name,
			}
		}
	}

	subTypes := make(map[string]discee)
	for _, sch := range doc.SchemasWithAllOf() {
		for _, ao := range sch.Schema.AllOf {
			if ao.Ref.String() != "" {
				if bt, ok := baseTypes[ao.Ref.String()]; ok {
					dce := discee{
						FieldName: bt.FieldName,
						// TODO: more trickery to allow for name customization
						FieldValue: swag.ToGoName(sch.Name),
						Ref:        sch.Ref,
						ParentRef:  ao.Ref,
						JSONName:   sch.Name,
						// TODO: more trickery to allow for name customization
						GoType: swag.ToGoName(sch.Name),
					}
					subTypes[sch.Ref.String()] = dce
					bt.Children = append(bt.Children, dce)
					baseTypes[ao.Ref.String()] = bt
				}
			}
		}
	}
	return &discInfo{Discriminators: baseTypes, Discriminated: subTypes}
}
