job "metermaid" {
    datacenters = ["us-west-2"]
    type = "service"

    constraint {
        attribute = "${meta.enclave}"
        value     = "shared"
    }

    meta {
        # Replace with version of metermaid to use
        VERSION = "v0.1.4"
    }

    group "primary" {
        // For testing purposes
        count = 1

        task "metermaid" {
            artifact {
                source = "https://github.com/euforia/metermaid/releases/download/${NOMAD_META_VERSION}/metermaid-linux.tgz"
            }

            driver = "raw_exec"
            config {
                command = "local/metermaid"
                args = [
                    "-bind-addr", "0.0.0.0:${NOMAD_PORT_default}",
                    "-adv-addr", "${NOMAD_ADDR_default}",
                    "-tags", "enclave=${meta.enclave},env=${meta.env}"
                ]
            }

            service {
                name = "metermaid"
                port = "default"
            }

            resources {
                cpu    = 100
                memory = 128
                network {
                    mbits = 1
                    port "default" {}
                }
            }
        }
    }
}