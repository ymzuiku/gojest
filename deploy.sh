v=v1.1.15
git tag $v
git push --tags
go install github.com/ymzuiku/gojest@$v
echo "done."