{{/*
Expand the name of the chart.
*/}}
{{- define "quibble.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "quibble.fullname" -}}
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
Chart label.
*/}}
{{- define "quibble.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels.
*/}}
{{- define "quibble.labels" -}}
helm.sh/chart: {{ include "quibble.chart" . }}
{{ include "quibble.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels — component is passed as a subcontext string.
*/}}
{{- define "quibble.selectorLabels" -}}
app.kubernetes.io/name: {{ include "quibble.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Component-scoped selector labels.
*/}}
{{- define "quibble.componentLabels" -}}
app.kubernetes.io/name: {{ include "quibble.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: {{ . }}
{{- end }}

{{/*
Service account name.
*/}}
{{- define "quibble.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "quibble.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Image tag — fall back to appVersion.
*/}}
{{- define "quibble.apiImage" -}}
{{- $tag := .Values.api.image.tag | default .Chart.AppVersion }}
{{- printf "%s:%s" .Values.api.image.repository $tag }}
{{- end }}

{{- define "quibble.webImage" -}}
{{- $tag := .Values.web.image.tag | default .Chart.AppVersion }}
{{- printf "%s:%s" .Values.web.image.repository $tag }}
{{- end }}
