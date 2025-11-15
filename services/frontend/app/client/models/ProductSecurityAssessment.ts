/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Issue } from './Issue';
export type ProductSecurityAssessment = {
    summary: string;
    key_issues: Array<Issue>;
    verdict: string;
    security_score: string;
    security_score_justification: string;
};

