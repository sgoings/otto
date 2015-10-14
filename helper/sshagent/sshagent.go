// Helpers for interacting with the local SSH Agent
package sshagent

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"reflect"
	"strings"

	execHelper "github.com/hashicorp/otto/helper/exec"
	"github.com/hashicorp/otto/ui"
  "github.com/mitchellh/go-homedir"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// HasKey determines if a given public key (provided as a string with the
// contents of a public key file), is loaded into the local SSH Agent.
func HasKey(publicKey string) (bool, error) {
	pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
	if err != nil {
		return false, fmt.Errorf("Error parsing provided public key: %s", err)
	}

	agentKeys, err := ListKeys()
	if err != nil {
		return false, err
	}

	for _, agentKey := range agentKeys {
		if reflect.DeepEqual(agentKey.Marshal(), pk.Marshal()) {
			return true, nil
		}
	}
	return false, nil
}

// ListKeys connects to the local SSH Agent and lists all the public keys
// loaded into it. It returns user friendly error message when it has trouble.
func ListKeys() ([]*agent.Key, error) {
	sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
	if sshAuthSock == "" {
		return nil, fmt.Errorf(
			"The SSH_AUTH_SOCK environment variable is not set, which normally\n" +
				"means that no SSH Agent is running.")
	}

	conn, err := net.Dial("unix", sshAuthSock)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to agent: %s", err)
	}
	defer conn.Close()

	agent := agent.NewClient(conn)
	loadedKeys, err := agent.List()
	if err != nil {
		return nil, fmt.Errorf("Error listing keys: %s", err)
	}
	return loadedKeys, err
}

// Add takes the path of a private key and runs ssh-add locally to add it to
// the agent. It needs a Ui to be able to interact with the user for the
// password prompt.
func Add(ui ui.Ui, privateKeyPath string) error {
	cmd := exec.Command("ssh-add", privateKeyPath)
	return execHelper.Run(ui, cmd)
}

func VerifyCreds(ui ui.Ui, publicKeyPath string) error {
  found, err := HasKey(publicKeyPath)
  if err != nil {
    return SshAgentError(err)
  }
  if !found {
    ok, _ := GuessAndLoadPrivateKey(
      ui, publicKeyPath)
    if ok {
      ui.Message(
        "A private key was found and loaded. Otto will now check\n" +
          "the SSH Agent again and continue if the correct key is loaded")

      found, err = HasKey(publicKeyPath)
      if err != nil {
        return SshAgentError(err)
      }
    }
  }

  if !found {
    return SshAgentError(fmt.Errorf(
      "You specified an SSH public key of: %q, but the private key from this\n"+
        "keypair is not loaded the SSH Agent. To load it, run:\n\n"+
        "  ssh-add [PATH_TO_PRIVATE_KEY]",
      publicKeyPath))
  }
  return nil
}

func SshAgentError(err error) error {
  return fmt.Errorf(
    "Otto uses your SSH Agent to authenticate with instances\n"+
      "but it could not verify that your SSH key is loaded into the agent.\n"+
      "The error message follows:\n\n%s", err)
}

// GuessAndLoadPrivateKey takes a path to a public key and determines if a
// private key exists by just stripping ".pub" from the end of it. if so,
// it attempts to load that key into the agent.
func GuessAndLoadPrivateKey(ui ui.Ui, pubKeyPath string) (bool, error) {
  fullPath, err := homedir.Expand(pubKeyPath)
  if err != nil {
    return false, err
  }
  if !strings.HasSuffix(fullPath, ".pub") {
    return false, fmt.Errorf("No .pub suffix, cannot guess path.")
  }
  privKeyGuess := strings.TrimSuffix(fullPath, ".pub")
  if _, err := os.Stat(privKeyGuess); os.IsNotExist(err) {
    return false, fmt.Errorf("No file at guessed path.")
  }

  ui.Header("Loading key into SSH Agent")
  ui.Message(fmt.Sprintf(
    "The key you provided (%s) was not found in your SSH Agent.", pubKeyPath))
  ui.Message(fmt.Sprintf(
    "However, Otto found a private key here: %s", privKeyGuess))
  ui.Message(fmt.Sprintf(
    "Automatically running 'ssh-add %s'.", privKeyGuess))
  ui.Message("If your SSH key has a passphrase, you will be prompted for it.")
  ui.Message("")

  if err := Add(ui, privKeyGuess); err != nil {
    return false, err
  }

  return true, nil
}
