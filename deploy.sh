v=v1.2.0
git tag $v
git push --tags
go install github.com/ymzuiku/gojest@$v
echo "done."