# Crypto_Stash

Crypto_Stash is an open-source CLI tool designed to securely manage your passwords and secrets through symmetric encryption. The program follows a simple and secure pre-loading flow, as illustrated in the diagram below:

![Pre-loading Flow](docs/base.png)

## How It Works

Here's a detailed explanation of the Crypto_Stash workflow:

1. **Key Generation:**
   - Upon initiation, Crypto_Stash generates a unique encryption key for you. This key is crucial for symmetric encryption and securing your secrets.

2. **Symmetric Encryption:**
   - Your secrets, such as passwords, are symmetrically encrypted using the generated key. This ensures that the confidentiality of your sensitive information is maintained.

3. **Secrets Storage:**
   - UUID reference to secrets are stored in a file named `secrets.json`. Each secret is uniquely identified by a randomly generated UUID, providing an additional layer of security.
4. **Secrets as file:**
    - Each encrypted secret is saved in a text file with the name being the UUID of the secret.
5. **Secure Access:**
   - Crypto_Stash allows you to access and manage your secrets securely through the CLI. With the generated key, you can decrypt and view your stored secrets as needed.

## Basic Pre-loading Flow

The diagram illustrates the foundational steps for Crypto_Stash to ensure a secure and efficient runtime.

### Getting Started

To use Crypto_Stash effectively, follow these basic steps:

1. **Compile the Application:**
   - Compile the application using `go build -o crypto_stash`.

2. **Generate a Secret Key:**
   - If no secret key is found, Crypto_Stash prompts you to generate one. This key is stored in `secret_key.txt` and is crucial for accessing your encrypted secrets.

3. **Manage Your Secrets:**
   - Use the CLI options to list, add, or retrieve individual secrets. The `secrets.json` file organizes your encrypted secrets for easy access.

### Notes

- Ensure the security of your secret key (`secret_key.txt`). Losing this key will result in permanent loss of access to your encrypted secrets.
- Exercise caution when sharing or distributing your compiled Crypto_Stash application to maintain the confidentiality of your secrets.

Crypto_Stash provides a secure and user-friendly solution for managing your sensitive information. Feel free to explore and customize the tool according to your specific needs.


> MADE By - Marlon Yepes Ceballos 
