v=v2.0.0
git tag $v
git push --tags
# go install 
go install github.com/ymzuiku/gojest@$v
echo "done."