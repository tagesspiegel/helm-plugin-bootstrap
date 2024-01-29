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

const pdbTemplate = `{{- if .Values.%[2]s.enabled }}
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: {{ include "%[1]s.fullname" . }}
  labels:
    {{- include "%[1]s.labels" . | nindent 4 }}
  {{- with .Values.%[2]s.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  {{- with .Values.%[2]s.maxUnavailable }}
  maxUnavailable: {{ . }}
  {{- end }}
  {{- with .Values.%[2]s.minAvailable }}
  minAvailable: {{ . }}
  {{- end }}
  selector:
    matchLabels:
	    {{- include "%[1]s.selectorLabels" . | nindent 6 }}
{{- end }}
`
const networkPolicyTemplate = `{{- if .Values.%[2]s.enabled }}
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
  {{- if .Values.%[2]s.ingress }}
    - Ingress
  {{- end }}
  {{- if .Values.%[2]s.egress }}
    - Egress
  {{- end }}
  {{- with .Values.%[2]s.ingress }}
  ingress:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  {{- with .Values.%[2]s.egress }}
  egress:
    {{- toYaml . | nindent 4 }}
  {{- end -}}
{{- end }}
`
const serviceMonitorTemplate = `{{- if and .Values.%[2]s.enabled .Values.%[2]s.serviceMonitor.enabled }}
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: {{ template "%[1]s.fullname" . }}
  {{- if .Values.%[2]s.serviceMonitor.namespace }}
  namespace: {{ .Values.%[2]s.serviceMonitor.namespace }}
  {{- end }}
  labels:
    {{- include "%[1]s.labels" . | nindent 4 }}
spec:
  endpoints:
    - port: {{ .Values.service.port }}
      path: {{ .Values.%[2]s.serviceMonitor.metricsPath }}
      {{- with .Values.%[2]s.serviceMonitor.interval }}
      interval: {{ . }}
      {{- end }}
      {{- with .Values.%[2]s.serviceMonitor.scrapeTimeout }}
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
