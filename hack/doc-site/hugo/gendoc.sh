#! /bin/bash
git clone https://github.com/alex-shpak/hugo-book themes/hugo-book
hugo server --config hugo.yaml,goswagger.yaml \
            --buildDrafts \
            --cleanDestinationDir \
            --minify \
            --printPathWarnings \
            --ignoreCache \
            --noBuildLock \
            --logLevel info \
            --source $(pwd)
