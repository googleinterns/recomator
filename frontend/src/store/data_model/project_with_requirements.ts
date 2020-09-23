export class Requirement {
  name: string;
  satisfied: boolean;
  errorMessage: string;

  constructor(name: string, satisfied: boolean, errorMessage: string) {
    this.name = name;
    this.satisfied = satisfied;
    this.errorMessage = errorMessage;
  }
}

export class ProjectRequirement {
  name: string;
  requirements: Requirement[];

  satisfiesRequirement(requirementName: string): boolean {
    for (const requirement of this.requirements) {
      if (requirement.name === requirementName) {
        return requirement.satisfied;
      }
    }

    return false;
  }

  hasRequirement(requirementName: string): boolean {
    for (const requirement of this.requirements) {
      if (requirement.name === requirementName) {
        return true;
      }
    }

    return false;
  }

  getErrorMessage(requirementName: string): string {
    for (const requirement of this.requirements) {
      if (requirement.name === requirementName) {
        return requirement.errorMessage;
      }
    }

    return "Unknown requirement";
  }

  constructor(project: string, requirements: Requirement[]) {
    this.name = project;
    this.requirements = new Array<Requirement>();

    for (const elt of requirements) {
      this.requirements.push(elt);
    }
  }
}
