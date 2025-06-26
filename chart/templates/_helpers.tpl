{{/*
Expand the name of the chart.
*/}}
{{- define "chart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "chart.fullname" -}}
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
{{- define "chart.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "chart.labels" -}}
helm.sh/chart: {{ include "chart.chart" . }}
{{ include "chart.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "chart.selectorLabels" -}}
app.kubernetes.io/name: {{ include "chart.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "chart.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "chart.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Get the Elasticsearch URL
*/}}
{{- define "chart.elasticsearchUrl" -}}
{{- if .Values.elasticsearch.enabled }}
{{- printf "https://%s-eck-elasticsearch-es-http:9200" .Release.Name }}
{{- else }}
{{- .Values.config.elasticsearchUrl }}
{{- end }}
{{- end }}

{{/*
Get the Elasticsearch username
*/}}
{{- define "chart.elasticsearchUsername" -}}
{{- if .Values.elasticsearch.enabled }}
{{- printf "elastic" }}
{{- else if .Values.config.elasticsearchUsername }}
{{- .Values.config.elasticsearchUsername }}
{{- end }}
{{- end }}

{{/*
Get elasticsearch auth secret name
*/}}
{{- define "chart.esPasswordSecretName" -}}
{{- if .Values.elasticsearch.enabled }}
{{- printf "%s-eck-elasticsearch-es-elastic-user" .Release.Name }}
{{- else if .Values.config.existingSecret.enabled }}
{{- .Values.config.existingSecret.secretName }}
{{- end }}
{{- end }}

{{/*
Get elasticsearch auth secret key
*/}}
{{- define "chart.esPasswordKey" -}}
{{- if .Values.elasticsearch.enabled }}
{{- printf "elastic" }}
{{- else if .Values.config.existingSecret.enabled }}
{{- .Values.config.existingSecret.passwordKey }}
{{- end }}
{{- end }}