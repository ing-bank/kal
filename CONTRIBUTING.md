# Contributing to KAL

`KAL` aims to overcome the problem of incomplete authorization listing in Kubernetes Clusters. This issue is most related to custom authorization endpoints, which sometimes does not inform an authorization response ([see more here](https://raesene.github.io/blog/2024/04/22/Fun-with-Kubernetes-Authz/)).

We're very much open to contributions but there are some things to keep in mind:

- Discuss the feature and implementation you want to add on Github before you write a PR for it. On disagreements, maintainer(s) will have the final word.
- Features need a somewhat general use case. If the use case is very niche it will be hard for us to consider maintaining it.
- If you’re going to add a feature, consider if you could help out in the maintenance of it.
- When issues or pull requests are not going to be resolved or merged, they should be closed as soon as possible. This is kinder than deciding this after a long period. Our issue tracker should reflect work to be done.

That said, there are many ways to contribute to KAL, including:

- Contribution to code
- Improving the documentation
- Reviewing merge requests
- Investigating bugs
- Reporting issues

<!-- Starting out with open source? See the guide [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/) and have a look at [our issues labelled *good first issue*](https://github.com/ing-bank/probatus/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22). -->


Planned improvements and changes are discussed in GitHub issues. Feel free to open a discussion.

Roadmap and future features are managed through GitHub projects.

## Standards

- Go latest version
- Formatter: [gofumpt](https://github.com/mvdan/gofumpt)
- Linter: [golangci-lint](https://github.com/golangci/golangci-lint)

### Branch structure

KAL follows the [Trunk Base Development]() structure.

- **master**: release code
- **feature/<your_feature>**: create based on **master** branch, contains new features and improvements
- **bugfix/<your_bugfix>**: create based on **master** branch, contains bug fixes

## Feature requests

If you have an improvement idea on how we can improve KAL, please check our [discussions]() and [projects]() to see if there are similar ideas or feature requests. If there is no similar idea, feel free to create a [new discussion]() with your idea. Use the **Feature request** template.

## Pull requests

Help out improving and fixing bugs in KAL by sending a Pull Request. Make sure that before contributing you create a fork of KAL and develop in your source version. After completing the development, open a [Pull Request](https://github.com/ing-bank/kal/pulls) following the **Pull Request** issue template.

Make sure no tests are broken and your code is linted:

```sh
make test
make lint
```

### Increase you PR approval

- Write tests
- Add documentation
- Write a [good commit message](https://www.conventionalcommits.org/).
