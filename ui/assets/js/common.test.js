import { ExpandVars, UrlParts, UrlToRepo } from "./common";

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

describe("UrlParts", () => {
    test("Generate parts from repo with default values", () => {
        const repo = {
            url: "https://www.github.com/YourOrganization/RepoOne.git",
            "url-pattern":
            {
                anchor: "#L{line}"
            }
        };
        const path = "test.txt"
        const line = null
        const rev = "main"
        expect(UrlParts(repo, path, line, rev)).toEqual(expect.objectContaining({
            url: "https://www.github.com/YourOrganization/RepoOne",
            rev: rev,
            path: path,
            anchor: "",
        }));
    });

    test("Generate parts from repo with default values and line", () => {
        const repo = {
            url: "https://www.github.com/YourOrganization/RepoOne.git",
            "url-pattern":
            {
                anchor: "#L{line}"
            }
        };
        const path = "test.txt"
        const line = "12"
        const rev = "main"
        expect(UrlParts(repo, path, line, rev)).toEqual(expect.objectContaining({
            url: "https://www.github.com/YourOrganization/RepoOne",
            rev: rev,
            path: path,
            anchor: "#L" + line,
        }));
    });

    test("Generate parts for ssh style repo with default values", () => {
        const repo = {
            url: "git@github.com:YourOrganization/RepoOne.git",
            "url-pattern":
            {
                anchor: "#L{line}"
            }
        };
        const path = "test.txt"
        const line = null
        const rev = "main"
        expect(UrlParts(repo, path, line, rev)).toEqual(expect.objectContaining({
            url: "//github.com/YourOrganization/RepoOne",
            rev: rev,
            path: path,
            anchor: "",
        }));
    });

    test("Generate parts for ssh bitbucket mercurial style repo", () => {
        const repo = {
            url: "ssh://hg@bitbucket.org/YourOrganization/RepoOne",
            "url-pattern":
            {
                "anchor" : "#{filename}-{line}"
            }
        };
        const path = "test.txt"
        const line = null
        const rev = "main"
        expect(UrlParts(repo, path, line, rev)).toEqual(expect.objectContaining({
            url: "//bitbucket.org/YourOrganization/RepoOne",
            path: path,
            anchor: "",
        }));
    });

    test("Generate parts for ssh bitbucket style repo with port", () => {
        const repo = {
            url: "ssh://git@bitbucket.org:7999/YourOrganization/RepoOne",
            "url-pattern":
            {
                "anchor" : "#{filename}-{line}"
            }
        };
        const path = "test.txt"
        const line = null
        const rev = "main"
        expect(UrlParts(repo, path, line, rev)).toEqual(expect.objectContaining({
            url: "//bitbucket.org:7999/YourOrganization/RepoOne",
            path: path,
            anchor: "",
        }));
    });

    test("Generate parts for ssh bitbucket server style repo", () => {
        const repo = {
            url: "ssh://git@bitbucket.internal.com:7999/YourOrganization/RepoOne",
            "url-pattern":
            {
                "base-url" : "{hostname}/projects/{project}/repos/{repo}/browse/{path}?at={rev}{anchor}",
                "anchor" : "#{line}",
            }
        };
        const path = "test.txt"
        const line = 10
        const rev = "main"
        expect(UrlParts(repo, path, line, rev)).toEqual(expect.objectContaining({
            hostname: "//bitbucket.internal.com",
            project: "YourOrganization",
            repo: "RepoOne",
            path: path,
            rev: rev,
            anchor: "#" + line,
        }));
    });
})

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
        const rev = "main"
        expect(UrlToRepo(repo, path, line, rev)).toBe(
            "https://www.github.com/YourOrganization/RepoOne/blob/main/test.txt"
        );
    });
});
