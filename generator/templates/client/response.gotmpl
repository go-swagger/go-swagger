{{- define "clientresponse" }}
// New{{ pascalize .Name }} creates a {{ pascalize .Name }} with default headers values
func New{{ pascalize .Name }}({{ if eq .Code -1 }}code int{{ end }}{{ if .Schema }}{{ if and (eq .Code -1) .Schema.IsStream }}, {{end}}{{ if .Schema.IsStream }}writer io.Writer{{ end }}{{ end }}) *{{ pascalize .Name }} {
  {{- if .Headers.HasSomeDefaults }}
  var (
  // initialize headers with default values
    {{- range .Headers }}
      {{- if .HasDefault }}
        {{ template "simpleschemaDefaultsvar" . }}
       {{- end }}
    {{- end }}
  )
    {{- range .Headers }}
      {{- if .HasDefault }}
        {{ template "simpleschemaDefaultsinit" . }}
      {{- end }}
    {{- end }}
  {{- end }}
  return &{{ pascalize .Name }}{
    {{- if eq .Code -1 }}
    _statusCode: code,
    {{- end }}
    {{ range .Headers }}
      {{- if .HasDefault }}
    {{ pascalize .Name}}: {{ if and (not .IsArray) (not .HasDiscriminator) (not .IsInterface) (not .IsStream) .IsNullable }}&{{ end }}{{ varname .ID }}Default,
      {{- end }}
    {{- end }}
    {{- if .Schema }}
      {{- if .Schema.IsStream }}
    Payload: writer,
      {{- end }}
    {{- end }}
    }
}

/* {{ pascalize .Name}} describes a response with status code {{ .Code }}, with default header values.

 {{ if .Description }}{{ blockcomment .Description }}{{else}}{{ pascalize .Name }} {{ humanize .Name }}{{end}}
 */
type {{ pascalize .Name }} struct {
  {{- if eq .Code -1 }}
  _statusCode int
  {{- end }}
  {{- range .Headers }}
    {{- if .Description }}

  /* {{ blockcomment .Description }}
     {{- if or .SwaggerFormat .Default }}
       {{ print "" }}
       {{- if .SwaggerFormat }}
     Format: {{ .SwaggerFormat }}
       {{- end }}
       {{- if .Default }}
     Default: {{ json .Default }}
       {{- end }}
     {{- end }}
  */
    {{- end }}
  {{ pascalize .Name }} {{ .GoType }}
  {{- end }}
  {{- if .Schema }}

  Payload {{ if and (not .Schema.IsBaseType) (not .Schema.IsInterface) .Schema.IsComplexObject (not .Schema.IsStream) }}*{{ end }}{{ if (not .Schema.IsStream) }}{{ .Schema.GoType }}{{ else }}io.Writer{{end}}
  {{- end }}
}
  {{- if eq .Code -1 }}

// Code gets the status code for the {{ humanize .Name }} response
func ({{ .ReceiverName }} *{{ pascalize .Name }}) Code() int {
  return {{ .ReceiverName }}._statusCode
}
  {{- end }}


func ({{ .ReceiverName }} *{{ pascalize .Name }}) Error() string {
	return fmt.Sprintf("[{{ upper .Method }} {{ .Path }}][%d] {{ if .Name }}{{ .Name }} {{ else }}unknown error {{ end }}{{ if .Schema }} %+v{{ end }}", {{ if eq .Code -1 }}{{ .ReceiverName }}._statusCode{{ else }}{{ .Code }}{{ end }}{{ if .Schema }}, o.Payload{{ end }})
}

  {{- if .Schema }}
func ({{ .ReceiverName }} *{{ pascalize .Name }}) GetPayload() {{ if and (not .Schema.IsBaseType) (not .Schema.IsInterface) .Schema.IsComplexObject (not .Schema.IsStream) }}*{{ end }}{{ if (not .Schema.IsStream) }}{{ .Schema.GoType }}{{ else }}io.Writer{{end}} {
	return o.Payload
}
  {{- end }}

func ({{ .ReceiverName }} *{{ pascalize .Name }}) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {
  {{- range .Headers }}

  // hydrates response header {{.Name}}
  hdr{{ pascalize .Name }} := response.GetHeader("{{ .Name }}")

  if hdr{{ pascalize .Name }} != "" {
    {{- if .Converter }}
  val{{ camelize .Name }}, err := {{ .Converter }}(hdr{{ pascalize .Name }})
  if err != nil {
    return errors.InvalidType({{ .Path }}, "header", "{{ .GoType }}", hdr{{ pascalize .Name }})
  }
  {{ .ReceiverName }}.{{ pascalize .Name }} = val{{ camelize .Name }}
    {{- else if .Child }}

  // binding header items for {{ .Name }}
  val{{ pascalize .Name }}, err := {{ .ReceiverName }}.bindHeader{{ pascalize .Name }}(hdr{{ pascalize .Name }}, formats)
  if err != nil {
    return err
  }

  {{ .ReceiverName }}.{{ pascalize .Name }} = val{{ pascalize .Name }}
    {{- else if .IsCustomFormatter }}
  val{{ camelize .Name }}, err := formats.Parse({{ printf "%q" .SwaggerFormat }}, hdr{{ pascalize .Name }})
  if err != nil {
    return errors.InvalidType({{ .Path }}, "header", "{{ .GoType }}", hdr{{ pascalize .Name }})
  }
      {{- if .IsNullable }}
  v := (val{{ camelize .Name }}.({{ .GoType }}))
  {{ .ReceiverName }}.{{ pascalize .Name }} = &v
      {{- else }}
  {{ .ReceiverName }}.{{ pascalize .Name }} = *(val{{ camelize .Name }}.(*{{ .GoType }}))
      {{- end }}
    {{- else }}
      {{- if eq .GoType "string" }}
  {{ .ReceiverName }}.{{ pascalize .Name }} = hdr{{ pascalize .Name }}
      {{- else }}
  {{ .ReceiverName }}.{{ pascalize .Name }} = {{ .GoType }}(hdr{{ pascalize .Name }})
      {{- end }}
    {{- end }}
  }
  {{-  end }}

  {{- if .Schema }}
    {{- if .Schema.IsBaseType }}

  // response payload as interface type
  payload, err := {{ toPackageName .ModelsPackage }}.Unmarshal{{ dropPackage .Schema.GoType }}{{ if .Schema.IsArray}}Slice{{ end }}(response.Body(), consumer)
  if err != nil {
    return err
  }
  {{ .ReceiverName }}.Payload = payload
    {{- else if .Schema.IsComplexObject }}

  {{ .ReceiverName }}.Payload = new({{ .Schema.GoType }})
    {{- end }}
    {{- if not .Schema.IsBaseType }}

  // response payload
  if err := consumer.Consume(response.Body(), {{ if not (or .Schema.IsComplexObject .Schema.IsStream) }}&{{ end}}{{ .ReceiverName }}.Payload); err != nil && err != io.EOF {
    return err
  }
    {{- end }}
  {{- end }}

  return nil
}
  {{- range .Headers }}
    {{- if .Child }}

// bindHeader{{ pascalize $.Name }} binds the response header {{ .Name }}
func ({{ .ReceiverName }} *{{ pascalize $.Name }}) bindHeader{{ pascalize .Name }}(hdr string, formats strfmt.Registry) ({{ .GoType }}, error) {
  {{ varname .Child.ValueExpression }}V := hdr

  {{ template "sliceclientheaderbinder" . }}

  return {{ varname .Child.ValueExpression }}C, nil
}
    {{- end }}
  {{- end }}
{{- end }}
// Code generated by go-swagger; DO NOT EDIT.


{{ if .Copyright -}}// {{ comment .Copyright -}}{{ end }}


package {{ .Package }}

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command


import (
  "io"
  "net/http"

  "github.com/go-openapi/errors"
  "github.com/go-openapi/runtime"
  "github.com/go-openapi/strfmt"
  "github.com/go-openapi/swag"
  "github.com/go-openapi/validate"

  {{ imports .DefaultImports }}
  {{ imports .Imports }}
)

// {{ pascalize .Name }}Reader is a Reader for the {{ pascalize .Name }} structure.
type {{ pascalize .Name }}Reader struct {
  formats strfmt.Registry
{{- if .HasStreamingResponse }}
  writer  io.Writer
{{- end }}
}

// ReadResponse reads a server response into the received {{ .ReceiverName }}.
func ({{ .ReceiverName }} *{{ pascalize .Name }}Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
  {{- if .Responses}}
  switch response.Code() {
  {{- end }}
  {{- range .Responses }}
    case {{ .Code }}:
      result := New{{ pascalize .Name }}({{ if .Schema }}{{ if .Schema.IsStream }}{{ $.ReceiverName }}.writer{{ end }}{{ end }})
      if err := result.readResponse(response, consumer, {{ $.ReceiverName }}.formats); err != nil {
        return nil, err
      }
      return {{ if .IsSuccess }}result, nil{{else}}nil, result{{ end }}
  {{- end }}
  {{- if .DefaultResponse }}
    {{- with .DefaultResponse }}
      {{- if $.Responses}}
    default:
      {{- end }}
      result := New{{ pascalize .Name }}(response.Code(){{ if .Schema }}{{ if .Schema.IsStream }}, {{ $.ReceiverName }}.writer{{ end }}{{ end }})
      if err := result.readResponse(response, consumer, {{ $.ReceiverName }}.formats); err != nil {
        return nil, err
      }
      if response.Code() / 100 == 2 {
        return result, nil
      }
      return nil, result
    {{- end }}
  {{- else }}
    {{- if $.Responses}}
    default:
    {{- end }}
      return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
  {{- end }}
  {{- if .Responses}}
  }
  {{- end }}
}

{{ range .Responses }}
  {{ template "clientresponse" . }}
{{ end }}
{{ if .DefaultResponse }}
  {{ template "clientresponse" .DefaultResponse }}
{{ end }}

{{ range .ExtraSchemas }}
/*{{ pascalize .Name }} {{ template "docstring" . }}
swagger:model {{ .Name }}
*/
  {{- template "schema" . }}
{{- end }}

{{- define "sliceclientheaderbinder" }}
 {{- if .IsArray }}
 var (
   {{ varname .Child.ValueExpression }}C {{ .GoType }}
 )
 // {{ .Child.ItemsDepth }}CollectionFormat: {{ printf "%q" .CollectionFormat }}
 {{ varname .Child.ValueExpression }}R := swag.SplitByFormat({{ varname .Child.ValueExpression }}V, {{ printf "%q" .CollectionFormat }})

 for {{ if or .Child.IsCustomFormatter .Child.Converter }}{{ .IndexVar }}{{ else }}_{{ end }}, {{ varname .Child.ValueExpression }}IV := range {{ varname .Child.ValueExpression }}R {
   {{ template "sliceclientheaderbinder" .Child }}
   {{ varname .Child.ValueExpression }}C = append({{ varname .Child.ValueExpression }}C, {{ varname .Child.ValueExpression }}IC) // roll-up {{ .Child.GoType }} into {{ .GoType }}
 }

 {{- else }}
   // convert split string to {{ .GoType }}
   {{- if .IsCustomFormatter }}
 val, err := formats.Parse({{ printf "%q" .SwaggerFormat }}, {{ varname .ValueExpression }}IV)
 if err != nil {
   return nil, errors.InvalidType({{ .Path }}, "header{{ .ItemsDepth }}", "{{ .GoType }}", {{ varname .ValueExpression }}IV)
 }
     {{- if .IsNullable }}
 {{ varname .ValueExpression }}IC := (&val).(*{{ .GoType }})
     {{- else }}
 {{ varname .ValueExpression }}IC := val.({{ .GoType }})
     {{- end }}
   {{- else if .Converter }}
 val, err := {{- print " "}}{{ .Converter }}({{ varname .ValueExpression }}IV)
 if err != nil {
   return nil, errors.InvalidType({{ .Path }}, "header{{ .ItemsDepth }}", "{{ .GoType }}", {{ varname .ValueExpression }}IV)
 }
     {{- if .IsNullable }}
 {{ varname .ValueExpression }}IC := &val
     {{- else }}
 {{ varname .ValueExpression }}IC := val
     {{- end }}
   {{- else }}
   {{ varname .ValueExpression }}IC :=
     {{- if eq .GoType "string" }}
       {{- print " " }}{{ varname .ValueExpression }}IV
     {{- else }}
       {{- print " " }}fmt.Sprintf("%v", {{ varname .ValueExpression }}IV)
     {{- end }} // string as {{ .GoType }}
   {{- end }}
 {{- end }}
{{- end }}
