$url = "http://localhost:8012/adduser"

$jsonPayload = @{
    name = "pupa"
    password = "123456"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri $url -Method Post -Body $jsonPayload -ContentType "application/json"

$response