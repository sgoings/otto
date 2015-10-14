provider "digitalocean" {
  token = "${var.do_token}"
}

resource "digitalocean_droplet" "deis" {
  count = "${var.instances}"
  image = "coreos-stable"
  name = "${var.prefix}-${count.index+1}"
  region = "${var.region}"
  size = "${var.size}"
  backups = "False"
  ipv6 = "False"
  private_networking = "True"
  ssh_keys = ["${var.ssh_keys}"]
  user_data = "${file("/Users/sethgoings/go/src/github.com/deis/deis/contrib/coreos/user-data")}"
}
