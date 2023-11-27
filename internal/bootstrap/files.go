package bootstrap

const (
	// PodDisruptionBudgetFileName is the name of the file that will be created in the templates folder
	PodDisruptionBudgetFileName = "pdb.yaml"
	// NetworkPolicyFileName is the name of the file that will be created in the templates folder
	NetworkPolicyFileName = "networkpolicy.yaml"
)

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

const networkPolicy = `{{- if .Values.networkPolicy.enabled }}
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

const pdbValuesYaml = `
pdb:
  enabled: true
  annotations: {}
  minAvailable: 1
  maxUnavailable: 0
`
const networkPolicyValuesYaml = `
networkPolicy:
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
