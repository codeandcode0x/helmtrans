{{/*
Expand the name of the chart.
*/}}
{{- define ".name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define ".fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define ".chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define ".labels" -}}
helm.sh/chart: {{ include ".chart" . }}
{{ include ".selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define ".selectorLabels" -}}
app.kubernetes.io/name: {{ include ".name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app: {{ include ".name" . }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define ".serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include ".fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}


{{/*
Convert octal to decimal (e.g 644 => 420). For file permission modes, many people are more familiar with octal notation.
However, due to yaml/json limitations, all the Kubernetes resources require file modes to be reported in decimal.
*/}}
{{- define ".fileModeOctalToDecimal" -}}
  {{- $digits := splitList "" (toString .) -}}

  {{/* Make sure there are exactly 3 digits */}}
  {{- if ne (len $digits) 3 -}}
    {{- fail (printf "File mode octal expects exactly 3 digits: %s" .) -}}
  {{- end -}}

  {{/* Go Templates do not support variable updating, so we simulate it using dictionaries */}}
  {{- $accumulator := dict "res" 0 -}}
  {{- range $idx, $digit := $digits -}}
    {{- $digitI := atoi $digit -}}

    {{/* atoi from sprig swallows conversion errors, so we double check to make sure it is a valid conversion */}}
    {{- if and (eq $digitI 0) (ne $digit "0") -}}
      {{- fail (printf "Digit %d of %s is not a number: %s" $idx . $digit) -}}
    {{- end -}}

    {{/*  Make sure each digit is less than 8 */}}
    {{- if ge $digitI 8 -}}
      {{- fail (printf "%s is not a valid octal digit" $digit) -}}
    {{- end -}}

    {{/* Since we don't have math.Pow, we hard code */}}
    {{- if eq $idx 0 -}}
      {{/* 8^2 */}}
      {{- $_ := set $accumulator "res" (add (index $accumulator "res") (mul $digitI 64)) -}}
    {{- else if eq $idx 1 -}}
      {{/* 8^1 */}}
      {{- $_ := set $accumulator "res" (add (index $accumulator "res") (mul $digitI 8)) -}}
    {{- else -}}
      {{/* 8^0 */}}
      {{- $_ := set $accumulator "res" (add (index $accumulator "res") (mul $digitI 1)) -}}
    {{- end -}}
  {{- end -}}
  {{- "res" | index $accumulator | toString | printf -}}
{{- end -}}

