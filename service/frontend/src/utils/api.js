import env from "../environment";

const get = async (url, options) =>
    await send(url, { method: "GET", ...options });

const post = async (url, body, options) =>
    await send(url, {
        method: "POST",
        body: JSON.stringify(body),
        headers: {
            "Content-Type": "application/json",
        },
        ...options,
    });

const postForm = async (url, formData, options) =>
    await send(url, {
        method: "POST",
        body: formData,
        ...options,
    });

const send = async (url, options) => {
    const dest = env.apiUrl + url;

    options.headers = {
        ...options.headers,
        Authorization:
            "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJleHAiOjE2NDM0OTEyMTcsInVpZCI6Ijc1QTkxNTU5LTRBRDQtNDM3NC05RTcwLTU2RDNDODIxRThGMyIsInVzZXJuYW1lIjoidGVzdDcifQ.aq2d6RT0T-KiWsCkUme8zM5Q8fs9wwwKNmY4pTGRTtMeE09YWEs27Sedpy5KX49dtM7IYgSw0uA5gsG_4-guXW5fGTv5Gva-XslTMzxX5G6vj0sVVYIM2Ty2KhrOT_GbOqgTa9ibutAvQGUdziqdjLGEfwQ5_UrKQRq1bU50wO9nWjPkNIpGxSnbAvLJxf5feYWB0RtPblWRNoXTAkjAm83UszQh4w_mxFK_BeC92PVcyA_yVE-eaZ6eIITvzMOwxf7nurNdJnMiFlgnn66ngMluFZBH8DW82GV4wpWmEHx336DUPW8AJYQqh0prTxYK_Nseb3LaQ0EHlrKSfpA2eQ",
    };

    try {
        const res = await fetch(dest, options);

        switch (res.status) {
            case 200:
                const data = await res.json();
                return {
                    ok: true,
                    data: data,
                };
            default:
                try {
                    const error = await res.json();
                    return {
                        ok: false,
                        error: error.message,
                    };
                } catch (e) {
                    console.error(
                        `Failed to read response from ${dest} (status: ${res.status})`,
                        e
                    );
                    return {
                        ok: false,
                        error: e.toString() || res.statusText,
                    };
                }
        }
    } catch (e) {
        console.error(
            `Failed to make a ${options.method} request to ${dest}`,
            e
        );
        return {
            ok: false,
            error: e.toString(),
        };
    }
};

export { get, post, postForm };
