# How to Contribute to Eiko
This is a guide on how to make a contribution to Eiko's code base

## Get the source
You need to clone the repository and have it in your GOPATH, for later golang compilation.
```bash
git clone https://github.com/eiko-team/eiko.git  $GOPATH/src/eiko
```

## Get Dependencies
You need to have GO>=1.10.4 installed and fetch dependencies.
```bash
cd $GOPATH/src/eiko
get get ./...
```

## Build the source
You should now be ready to type `make` from the root of your `git` source directory.
Here are some helpful variations:

* `make up`:    Compile new sources and run the server in a docker container. You can now open [127.0.0.1](http://127.0.0.1) on your browser.
* `make clean`: Delete all compiled files.

## What to change ?
One way to win big with the open-source community is to look at the
[issues page](https://github.com/eiko-team/eiko/issues) and see if there are any issues that you can fix quickly, or if anything catches your eye. There is also some [styling issues](https://app.codacy.com/manual/tomMoulard/eiko/dashboard) that you can checkout.

Or you can "scratch your own itch", i.e. address an issue you have with Eiko. By either, filling a [new issue](https://github.com/eiko-team/eiko/issues/new) or by solving it first hand by making a pull request.

You could also think of something you wish Git could do, and make it do that thing! The only concern I would have with this approach is whether or not that feature is something the community also wants. If this excites you though, go for it! Don't be afraid to open a pull request.

## Change
Make your changes, and don't forget comments.

### Coding conventions
Start reading our code and you'll get the hang of it.
We use Codacity to manage code quality and coding conventions, and Go Report Card for golang specific coding conventions. Thus, you can go check [eiko's codacity](https://app.codacy.com/manual/tomMoulard/eiko/dashboard) to have a deep look at issues.

### README
Update the README.md with details of changes to the interface, this includes new environment variables, exposed ports, useful file locations and container parameters.

## Test Your Changes
After you make your changes, it is important that you test your changes. Manual testing is important, but checking and extending the existing test suite is even more important. You want to run the functional tests to see if you broke something else during your change, and you want to extend the functional tests to be sure no one breaks your feature in the future.

### Functional Tests
Using golang test packages, every package must have an associated package named `PACKAGE_test` with test functions inside(coverage will be checked). If you want an example on have to build a package and it's test package, take a look at [this package](https://github.com/eiko-team/eiko/tree/master/src).

To run tests, use `make` or `make test`. There will be a test summary with coverage associated.
For a better coverage visualization, you can use `make cover` after `make test` and then, open your browser with the file `test.html`.

### Polish Your Commits
Before submitting your patch, be sure to read the [codding guidelines](https://github.com/eiko-team/eiko/blob/master/.github/CODE_OF_CONDUCT.md) and check your code to match as best you can. This can be a lot of efforts, but it saves time during review to avoid style issues.

When preparing your patch, it is important to put yourself in the shoes of the Eiko community. Accepting a patch requires more justification than approving a pull request from someone on your team. The community has a stable product and is responsible for keeping it stable. If you introduce a bug, then they cannot count on you being around to fix it. When you decided to start work on a new feature, they were not part of the design discussion and may not even believe the feature is worth introducing.

Questions to answer in your patch message (and commit messages) may include:
* Why is this patch necessary?
* How does the current behavior cause pain for users?
* What kinds of repositories are necessary for noticing a difference?
* What design options did you consider before writing this version? Do you have links to
  code for those alternate designs?
* Is this a performance fix? Provide clear performance numbers for various well-known repos.

Here are some other tips that we use when cleaning up our commits:
* Make sure the commits are signed off using `git commit (-s|--signoff)`. See
  [SubmittingPatches](https://help.github.com/en/articles/signing-commits) for more details about what this sign-off means.
* Your commit titles should match the "area: change description" format. Rules of thumb:
    * Choose "&lt;area>: " prefix appropriately.
    * Keep the description short and to the point.
    * The word that follows the "&lt;area>: " prefix is not capitalized.
    * Do not include a full-stop at the end of the title.
    * Read a few commit messages (using `git log origin/master`, for instance) to become acquainted with the preferred commit message style.

## Submit Your Patch
Eiko [accepts pull requests on GitHub](https://github.com/eiko-team/eiko/pulls).
