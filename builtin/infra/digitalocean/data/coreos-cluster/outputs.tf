output "hosts" {
  value = "${join(", ", digitalocean_droplet.deis.*.ipv4_address)}"
}

output "ip" {
  value = "${digitalocean_droplet.deis.0.ipv4_address}"
}
