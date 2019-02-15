job "metermaid" {
    type = "system"

    meta {
        VERSION = ""
    }

    group "primary" {
        task "metermaid" {
            driver = "raw_exec"
            config {

            }
        }
    }
}