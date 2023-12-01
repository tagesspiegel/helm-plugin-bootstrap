package bootstrap

const (
	// PodDisruptionBudgetFileName is the name of the file that will be created in the templates folder
	PodDisruptionBudgetFileName = "pdb.yaml"
	// NetworkPolicyFileName is the name of the file that will be created in the templates folder
	NetworkPolicyFileName = "networkpolicy.yaml"
	// ServiceMonitorFileName is the name of the file that will be created in the templates folder
	ServiceMonitorFileName = "servicemonitor.yaml"
)

// manifest templates

const pdbTemplate = `{{- if .Values.pdb.enabled }}
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: {{ include "%[1]s.fullname" . }}
  labels:
    {{- include "%[1]s.labels" . | nindent 4 }}
  {{- with .Values.pdb.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  {{- with .Values.pdb.maxUnavailable }}
  maxUnavailable: {{ . }}
  {{- end }}
  {{- with .Values.pdb.minAvailable }}
  minAvailable: {{ . }}
  {{- end }}
  selector:
    matchLabels:
	{{- include "%[1]s.selectorLabels" . | nindent 6 }}
{{- end }}
`
const networkPolicyTemplate = `{{- if .Values.networkPolicy.enabled }}
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: {{ include "%[1]s.fullname" . }}
  labels:
    {{- include "%[1]s.labels" . | nindent 4 }}
spec:
  podSelector:
    matchLabels:
      {{- include "%[1]s.selectorLabels" . | nindent 6 }}
  policyTypes:
  {{- if .Values.networkPolicy.ingress }}
    - Ingress
  {{- end }}
  {{- if .Values.networkPolicy.egress }}
    - Egress
  {{- end }}
  {{- with .Values.networkPolicy.ingress }}
  ingress:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  {{- with .Values.networkPolicy.egress }}
  egress:
    {{- toYaml . | nindent 4 }}
  {{- end -}}
{{- end }}
`
const serviceMonitorTemplate = `{{- if and .Values.metrics.enabled .Values.metrics.serviceMonitor.enabled }}
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: {{ template "%[1]s.fullname" . }}
  {{- if .Values.metrics.serviceMonitor.namespace }}
  namespace: {{ .Values.metrics.serviceMonitor.namespace }}
  {{- end }}
  labels:
    {{- include "%[1]s.labels" . | nindent 4 }}
spec:
  endpoints:
    - port: {{ .Values.service.port }}
      path: {{ .Values.metrics.serviceMonitor.metricsPath }}
      {{- with .Values.metrics.serviceMonitor.interval }}
      interval: {{ . }}
      {{- end }}
      {{- with .Values.metrics.serviceMonitor.scrapeTimeout }}
      scrapeTimeout: {{ . }}
      {{- end }}
  selector:
    matchLabels:
      {{- include "%[1]s.selectorLabels" . | nindent 6 }}
{{- end }}
`

// values.yaml configurations

const pdbValuesYaml = `
%[1]s:
  enabled: true
  annotations: {}
  minAvailable: 1
  maxUnavailable: 0
`
const networkPolicyValuesYaml = `
%[1]s:
  enabled: false
  ingress: []
    # - from:
    #   - ipBlock:
    #       cidr: 10.0.0.0/24
    #       except:
    #         - 10.0.0.128/25
    #   - namespaceSelector:
    #       matchLabels:
    #         kubernetes.io/metadata.name: frontend
    #   - podSelector:
    #       matchLabels:
    #         app.kubernetes.io/name: frontend
    #   ports:
    #     - protocol: TCP
    #       port: 80
  egress: []
    # - to:
    #     - ipBlock:
    #         cidr: 10.0.0.0/24
    #   ports:
    #     - protocol: UDP
    #       port: 53
`
const serviceMonitorValuesYaml = `
%[1]s:
  enabled: false
  serviceMonitor:
    enabled: false
    metricsPath: /metrics
    namespace: ""
    interval: ""
    scrapeTimeout: ""
`
