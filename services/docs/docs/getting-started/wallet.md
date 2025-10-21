---

title: Wallet
---
# Lethean Wallet

!!! info

    === "Windows"

        [Download Windows CLI](https://github.com/letheanVPN/blockchain-iz/releases/latest/download/windows.tar){ .md-button }
        [Blockchain Export](https://seed.lethean.io/blockchain.raw){ .md-button }

        Remote Hosts: `seed.lethean.io`, `nodes.hashvault.pro`

    === "MacOS"

        [Download macOS CLI](https://github.com/letheanVPN/blockchain-iz/releases/latest/download/lethean-cli-macos.zip){ .md-button }
        [Blockchain Export](https://seed.lethean.io/blockchain.raw){ .md-button }

        Remote Hosts: `seed.lethean.io`, `nodes.hashvault.pro`
      

    === "Linux"

        [Download Linux CLI](https://github.com/letheanVPN/blockchain-iz/releases/latest/download/linux.tar){ .md-button }
        [Blockchain Export](https://seed.lethean.io/blockchain.raw){ .md-button }

        Remote Hosts: `seed.lethean.io`, `nodes.hashvault.pro`

## Data Location



=== "Windows"
    
    ```shell
    %USERPROFILE%\\Lethean\\wallets 
    ```

=== "MacOS"

    ``` shell
    $HOME/Lethean/wallets
    ```

=== "Linux"

    ``` shell
    $HOME/Lethean/wallets
    ```



## Using the CLI

!!! example

    === "Windows"

        1. Press the Windows key
        2. type `cmd.exe` + Press Enter
        3. change directory to Lethean user data `cd %USERPROFILE%\\Lethean\\`

    === "MacOS"

        1. Press the `CMD` + `SPACE` 
        2. type `Terminal` + Press Enter
        3. change directory to Lethean user data `cd $HOME/Lethean`

    === "Linux"

        1. Open your preferred Terminal
        2. change directory to Lethean user data `cd $HOME/Lethean`
        

### New Wallet 


=== "Windows"
    
    ```shell
    cd %USERPROFILE%\\Lethean\\wallets && ..\\cli\\lethean-wallet-cli.exe --daemon-host=seed.lethean.io --generate-new-wallet=wallet
    ```

=== "MacOS"

    ``` shell
    cd $HOME/Lethean/wallets && ../cli/lethean-wallet-cli --daemon-host=seed.lethean.io --generate-new-wallet=wallet
    ```

=== "Linux"

    ``` shell
    cd $HOME/Lethean/wallets && ../cli/lethean-wallet-cli --daemon-host=seed.lethean.io --generate-new-wallet=wallet
    ```

### Restore Wallet from Seed


=== "Windows"
    
    ```shell
    cd %USERPROFILE%\\Lethean\\wallets && ..\\cli\\lethean-wallet-cli.exe --daemon-host=seed.lethean.io -restore-deterministic-wallet --generate-new-wallet=wallet
    ```

=== "MacOS"

    ``` shell
    cd $HOME/Lethean/wallets && ../cli/lethean-wallet-cli --daemon-host=seed.lethean.io -restore-deterministic-wallet --generate-new-wallet=wallet
    ```

=== "Linux"

    ``` shell
    cd $HOME/Lethean/wallets && ../cli/lethean-wallet-cli --daemon-host=seed.lethean.io -restore-deterministic-wallet --generate-new-wallet=wallet
    ```


### Restore Wallet From Keys


=== "Windows"
    
    ```shell
    cd %USERPROFILE%\\Lethean\\wallets && ..\\cli\\lethean-wallet-cli.exe --daemon-host=seed.lethean.io --generate-new-keys=wallet
    ```

=== "MacOS"

    ``` shell
    cd $HOME/Lethean/wallets && ../cli/lethean-wallet-cli --daemon-host=seed.lethean.io --generate-new-keys=wallet
    ```

=== "Linux"

    ``` shell
    cd $HOME/Lethean/wallets && ../cli/lethean-wallet-cli --daemon-host=seed.lethean.io --generate-new-keys=wallet
    ```



### Open Wallet

=== "Windows"
    
    ```shell
    cd %USERPROFILE%\\Lethean\\wallets &&  ..\\cli\\lethean-wallet-cli --daemon-host=seed.lethean.io --wallet-file=wallet 
    ```

=== "MacOS"

    ``` shell
    cd $HOME/Lethean/wallets && ../cli/lethean-wallet-cli --daemon-host=seed.lethean.io --wallet-file=wallet
    ```

=== "Linux"

    ``` shell
    cd $HOME/Lethean/wallets && ../cli/lethean-wallet-cli --daemon-host=seed.lethean.io --wallet-file=wallet
    ```

