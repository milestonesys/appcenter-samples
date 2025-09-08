$ErrorActionPreference = "Stop"

Try {
    $Server = Read-Host -Prompt 'Input hostname of Management Server'
    $User = Read-Host -Prompt 'Input username of basic user in the administrator role'
    $Pass = Read-Host -MaskInput -Prompt 'Input password of basic user'

    $Endpoint = "https://" + $Server + ":/IDP/connect/token"

    $Response = Invoke-WebRequest -AllowUnencryptedAuthentication -Uri $Endpoint -Method Post -Body @{
        grant_type = "password"
        username = $User
        password = $Pass
        client_id = "GrantValidatorClient"
    }
    $BUF_AccessToken = ($Response.Content | ConvertFrom-Json).access_token

    $Response = Invoke-WebRequest -AllowUnencryptedAuthentication -Uri $Endpoint -Method Post -Body @{
        grant_type = "windows_auth"
        scope = "write:client"
        client_id = "winauthclient"
    } -UseDefaultCredentials
    $CCF_AccessToken = ($Response.Content | ConvertFrom-Json).access_token

    Set-Clipboard ""
    Set-Clipboard -Append -Value "kubectl create secret generic app-registration-buf-token --from-literal='token=$BUF_AccessToken' --dry-run=client -o yaml | kubectl apply -f -"
    Set-Clipboard -Append -Value "kubectl create secret generic app-registration-ccf-token --from-literal='token=$CCF_AccessToken' --dry-run=client -o yaml | kubectl apply -f -"

    Write-Output "Your clipboard now contains the commands for creating the secrets necessary for using the basic user flow and client credentials flow"
}
Catch {
    Write-Output $_.Exception
}

Write-Host "Press any key to continue"
$void = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")