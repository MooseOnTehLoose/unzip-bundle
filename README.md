# unzip-bundle
This tool will unzip a DC/OS bundle using pure go and no OS dependencies.

By default the tool will create a new folder in the same directory as the bundle and unzip its contents there. You can also specify where to drop the bundle files. 

Usage: unzip-bundle <source dc/os diag bundle> <optional destination path>
