	mkdir ~/.meshery
	config='{"contexts":{"local":{"endpoint":"http://localhost:9081","token":"Default","platform":"kubernetes","adapters":[],"channel":"stable","version":"latest"}},"current-context":"local","tokens":[{"location":"auth.json","name":"Default"}]}'

	echo $config > ~/.meshery/config.yaml

	cat ~/.meshery/config.yaml
	echo '{ "meshery-provider": "Meshery", "token": null }' | jq -c '.token = "'$provider_token'"' > ~/.meshery/auth.json
	cat ~/.meshery/auth.json	