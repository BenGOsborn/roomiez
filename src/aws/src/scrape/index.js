exports.handler = async (event, context) => {};

(async () => {
    if (process.env.ENV !== "production") {
        require("dotenv").config({ path: "../../.env" });
    }

    const AWS = require("aws-sdk");
    const secretsManager = new AWS.SecretsManager({ region: "ap-southeast-2" });
    const puppeteer = require("puppeteer-core");

    // Load env
    const secretsArn = process.env.SECRETS_ARN;

    const secretsString = (await secretsManager.getSecretValue({ SecretId: secretsArn }).promise()).SecretString;
    if (!secretsString) throw Error("could not retrieve secrets");

    const secrets = JSON.parse(secretsString);

    // Initialize browser
    let browser = await puppeteer.connect({ browserWSEndpoint: secrets["BROWSER_ENDPOINT"] });

    const page = await browser.newPage();
    page.setDefaultNavigationTimeout(2 * 60 * 1000);

    await page.goto("https://www.facebook.com/groups/2280085492006745");
    await page.waitForSelector('[role="feed"]');

    // **** If it works we will loop through this by scrolling down N times
    const html = await page.evaluate(async () => {
        const out = [];

        for (let i = 0; i < 2; i++) {
            window.scrollTo(0, document.body.scrollHeight);
            await new Promise((res) => setTimeout(res, 5000));
        }

        const elements = Array.from(document.querySelector('[role="feed"]').children);
        return elements.length;

        for (const elem of elements) {
            const msg = elem.querySelector('[data-ad-comet-preview="message"]');
            const more = elem.querySelector('[role="button"]');

            if (msg) {
                if (more) {
                    more.click();
                    await new Promise((res) => setTimeout(res, 100));
                }

                out.push(msg.textContent);
            }
        }

        return out;
    });

    console.log(html);

    await browser?.close();
})();
