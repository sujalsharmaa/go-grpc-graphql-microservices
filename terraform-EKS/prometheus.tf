resource "helm_release" "prometheus_helm" {
  name             = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  version          = "62.3.1"
  namespace        = "prometheus"
  create_namespace = true
  timeout          = 2000

  set {
    name  = "podSecurityPolicy.enabled"
    value = "true"
  }

  set {
    name  = "server.persistentVolume.enabled"
    value = "true"
  }

  set {
    name  = "grafana.service.type"
    value = "LoadBalancer"
  }

  set {
    name  = "prometheus.service.type"
    value = "LoadBalancer"
  }

  # Alertmanager SMTP settings (as a YAML block override)
  values = <<EOT
alertmanager:
  config:
    global:
      smtp_smarthost: "smtp.gmail.com:587"
      smtp_from: "techsharma53@gmail.com"
      smtp_auth_username: "techsharma53@gmail.com"
      smtp_auth_password: "csvdmvyncmjovwwx"
      smtp_require_tls: true
EOT

  depends_on = [helm_release.argocd]
}
