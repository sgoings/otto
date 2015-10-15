variable "instances" {
  default = "3"
}

variable "prefix" {
  default = "deis"
}

variable "region" {
  default = "nyc1"
}

variable "size" {
  default = "8gb"
}

variable "ssh_keys" {
  description = "The ssh fingerprint of the ssh key you'll be using"
}

variable "do_token" {
  description = "Your Digital Ocean auth token"
}

variable "user" {
  default = "sgoings"
}
