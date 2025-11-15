/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ReportContent } from './ReportContent';
export type ReportGetResponse = {
    id: string;
    query: string;
    status: ReportGetResponse.status;
    created_at: string;
    content?: ReportContent;
};
export namespace ReportGetResponse {
    export enum status {
        UNKNOWN = 'unknown',
        PENDING = 'pending',
        READY = 'ready',
        FAILED = 'failed',
    }
}

