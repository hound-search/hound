import { ExpandVars, UrlToRepo } from "./common";

describe("ExpandVars", () => {
    test("Replaces template variables with their values", () => {
        const template = "I am trying to {verb} my {noun}";
        const values = { verb: "wash", noun: "dishes" };
        expect(ExpandVars(template, values)).toBe(
            "I am trying to wash my dishes"
        );
    });

    test("Doesn't replace unlisted variables", () => {
        const template = "Get the {expletive} out of my {noun}";
        const values1 = { noun: "stamp collection" };

        expect(ExpandVars(template, values1)).toBe(
            "Get the {expletive} out of my stamp collection"
        );
        expect(ExpandVars(template, {})).toBe(template);
    });
});

describe("UrlToRepo", () => {
    test("Generate url from repo with default values", () => {
        const repo = {
            url: "https://www.github.com/YourOrganization/RepoOne.git",
            "url-pattern":
            {
                "base-url": "{url}/blob/{rev}/{path}{anchor}",
                anchor: "#L{line}"
            }
        };
        const path = "test.txt"
        const line = null
        const rev = "master"
        expect(UrlToRepo(repo, path, line, rev)).toBe(
            "https://www.github.com/YourOrganization/RepoOne/blob/master/test.txt"
        );
    });

    test("Generate url from repo with default values and line", () => {
        const repo = {
            url: "https://www.github.com/YourOrganization/RepoOne.git",
            "url-pattern":
            {
                "base-url": "{url}/blob/{rev}/{path}{anchor}",
                anchor: "#L{line}"
            }
        };
        const path = "test.txt"
        const line = "12"
        const rev = "master"
        expect(UrlToRepo(repo, path, line, rev)).toBe(
            "https://www.github.com/YourOrganization/RepoOne/blob/master/test.txt#L12"
        );
    });

    test("Generate url from for ssh style repo with default values", () => {
        const repo = {
            url: "git@github.com:YourOrganization/RepoOne.git",
            "url-pattern":
            {
                "base-url": "{url}/blob/{rev}/{path}{anchor}",
                anchor: "#L{line}"
            }
        };
        const path = "test.txt"
        const line = null
        const rev = "master"
        expect(UrlToRepo(repo, path, line, rev)).toBe(
            "//github.com/YourOrganization/RepoOne/blob/master/test.txt"
        );
    });

    test("Generate url from for ssh bitbucket mercurial style repo", () => {
        const repo = {
            url: "ssh://hg@bitbucket.org/YourOrganization/RepoOne",
            "url-pattern":
            {
                "base-url" : "{url}/src/master/{path}{anchor}",
                "anchor" : "#{filename}-{line}"
            }
        };
        const path = "test.txt"
        const line = null
        const rev = "master"
        expect(UrlToRepo(repo, path, line, rev)).toBe(
            "//bitbucket.org/YourOrganization/RepoOne/src/master/test.txt"
        );
    });
});
