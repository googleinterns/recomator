import { ProjectRequirement } from "./data_model/project_with_requirement";

export interface IRequirementsStoreState {
  projects: ProjectRequirement[];
  progress: null | number;
  errorCode: undefined | number;
  errorMessage: undefined | string;
  requestId: string;
  display: boolean;
}

export function requirementsStoreStateFactory(): IRequirementsStoreState {
  return {
    projects: [],
    progress: null,
    errorCode: undefined,
    errorMessage: undefined,
    requestId: "",
    display: false
  };
}
