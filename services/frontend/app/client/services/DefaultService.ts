/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ReportGetResponse } from '../models/ReportGetResponse';
import type { ReportStartResponse } from '../models/ReportStartResponse';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class DefaultService {
    /**
     * Start report generation
     * Starts async report generation and returns the report id.
     * @param requestBody
     * @returns ReportStartResponse Report creation accepted
     * @throws ApiError
     */
    public static createReport(
        requestBody: {
            /**
             * Query parameter for report generation
             */
            query: string;
        },
    ): CancelablePromise<ReportStartResponse> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/report',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                500: `Internal server error`,
            },
        });
    }
    /**
     * Get generated report
     * Returns the report object by id.
     * @param id
     * @returns ReportGetResponse A report object
     * @throws ApiError
     */
    public static getReportById(
        id: string,
    ): CancelablePromise<ReportGetResponse> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/report/{id}',
            path: {
                'id': id,
            },
            errors: {
                400: `Invalid report ID`,
                404: `Report not found`,
                500: `Internal server error`,
            },
        });
    }
}
