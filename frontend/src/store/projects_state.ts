import { Project } from "./data_model/project";

export interface IProjectsStoreState {
  projects: Project[];
  projectsSelected: Project[];
  errorCode: number | undefined;
  errorMessage: string | undefined;
  loading: boolean;
  loaded: boolean;
}

export function projectsStoreStateFactory(): IProjectsStoreState {
  return {
    projects: [],
    projectsSelected: [],
    errorCode: undefined,
    errorMessage: undefined,
    loading: false,
    loaded: false
  };
}
