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
