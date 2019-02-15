job "metermaid" {
    type = "system"

    meta {
        VERSION = ""
    }

    group "primary" {
        task "metermaid" {
            artifact {
                source = "https://github.com/euforia/metermaid/metermaid-linux.tgz"
            }

            driver = "raw_exec"
            config {
                command = "local/metermaid"
                args = [
                    "-bind-addr", "",
                    "", "${NOMAD_HOST_ADDR}",
                ]
            }
        }
    }
}