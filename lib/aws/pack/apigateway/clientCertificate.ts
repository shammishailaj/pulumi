// *** WARNING: this file was generated by the Coconut IDL Compiler (CIDLC).  ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as coconut from "@coconut/coconut";

export class ClientCertificate extends coconut.Resource implements ClientCertificateArgs {
    public readonly name: string;
    public description?: string;

    constructor(name: string, args: ClientCertificateArgs) {
        super();
        if (name === undefined) {
            throw new Error("Missing required resource name");
        }
        this.name = name;
        this.description = args.description;
    }
}

export interface ClientCertificateArgs {
    description?: string;
}

