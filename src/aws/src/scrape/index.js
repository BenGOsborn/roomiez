const chromium = require("chrome-aws-lambda");

exports.handler = async (event, context, callback) => {
    // Load env
    if (process.env.ENV !== "production") {
        require("dotenv").config({ path: "../../.env" });
    }

    const AWS = require("aws-sdk");
    const secretsManager = new AWS.SecretsManager({ region: "ap-southeast-2" });
    const puppeteer = require("puppeteer-core");
    const fs = require("fs");

    // Load env
    const secretsArn = process.env.SECRETS_ARN;

    const secretsString = (await secretsManager.getSecretValue({ SecretId: secretsArn }).promise()).SecretString;
    if (!secretsString) throw Error("could not retrieve secrets");

    const secrets = JSON.parse(secretsString);

    // Scrape
    let result = null;
    let browser = null;

    try {
        browser = await chromium.puppeteer.launch({
            args: chromium.args,
            defaultViewport: chromium.defaultViewport,
            executablePath: await chromium.executablePath,
            headless: chromium.headless,
            ignoreHTTPSErrors: true,
        });

        let page = await browser.newPage();
        page.setDefaultNavigationTimeout(2 * 60 * 1000);

        // Login
        await page.goto("https://www.facebook.com/login");
        await page.waitFor(3000);

        await page.type("[id=email]", secrets["FB_EMAIL"]);
        await page.type("[id=pass]", secrets["FB_PASS"]);
        await page.click("[type=submit]");
        await page.waitFor(3000);

        // Scrape data from group
        await page.goto("https://www.facebook.com/groups/2280085492006745");
        await page.waitFor(3000);

        for (let i = 0; i < 1; i++) await page.keyboard.press("PageDown");
        await page.waitFor(3000);

        // **** This code mostly needs to be converted to puppeteer code instead (especially for the URL scraping)
        const out = await page.evaluate(async () => {
            const out = [];

            for (const elem of Array.from(document.querySelector("[role=feed]").children)) {
                const msg = elem.querySelector("[data-ad-comet-preview=message]");

                // Get the post URL
                const url = await new Promise((res) => {
                    elem.querySelector('[aria-label="Send this to friends or post it on your Timeline."]')?.click();

                    const shareOptionContainer = elem.querySelector("[role=dialog]");
                    const shareOptions = shareOptionContainer.querySelectorAll("[role=button]");
                    const copyLink = shareOptions[shareOptions.length - 1];

                    copyLink.addEventListener("click", () => {
                        navigator.clipboard.readText().then((clipboardText) => {
                            res(clipboardText);
                        });
                    });

                    copyLink.click();
                });

                // Get the text
                const more = elem.querySelector("[role=button]");

                if (msg && url) {
                    if (more) {
                        more.click();
                        await new Promise((res) => setTimeout(res, 100));
                    }

                    out.push({ post: msg.textContent, url });
                }
            }

            return out;
        });

        return out;
    } catch (error) {
        return callback(error);
    } finally {
        await browser?.close();
    }
};
