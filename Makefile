release-with-gpg-pass:
	@dagger call -m gopkg with-git-gpg-config \
		--gpg-key=env://GPG_KEY \
		--gpg-key-id=env://GPG_KEY_ID \
		--gpg-password=env://GPG_PASSPHRASE \
		--git-author-name "semantic-release-bot" \
		--git-author-email "shukomango@gmail.com" \
		release \
		--source=. \
		--dry-run=false \
		--token=env://GITHUB_TOKEN

release-with-gpg-no-pass:
	@dagger call -m gopkg with-git-gpg-config \
		--gpg-key=env://GPG_KEY \
		--gpg-key-id=env://GPG_KEY_ID \
		--git-author-name "semantic-release-bot" \
		--git-author-email "shukomango@gmail.com" \
		release \
		--source=. \
		--dry-run=false \
		--token=env://GITHUB_TOKEN



release:
	@dagger call -m gopkg release \
		--source=. \
		--dry-run=false \
		--token=env://GITHUB_TOKEN


