import { ProjectRequirement } from "./data_model/project_with_requirements";

export interface IRequirementsStoreState {
  projects: ProjectRequirement[];
  progress: null | number;
  requestId: string;
  display: boolean;
}

export function requirementsStoreStateFactory(): IRequirementsStoreState {
  return {
    projects: [],
    progress: null,
    requestId: "",
    display: false
  };
}
