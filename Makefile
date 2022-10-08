


test:
	docker build -f test/Dockerfile -t fact_app_test .
	docker run -v ${PWD}/cover.out:/testdir/cover.out -e GIT_URL='' fact_app_test

.PHONY: test