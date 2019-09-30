#!/bin/bash

tar --exclude='./node_modules' --exclude='./package' --exclude='./dist' -zcvf extension.tar.gz * .babelrc .eslintrc
