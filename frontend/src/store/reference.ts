// Class that allows to pass by reference
export class ReferenceWrapper<T> {
  private variableWrapper: T;

  constructor(initialValue: T) {
    this.variableWrapper = initialValue;
  }

  public getValue(): T {
    return this.variableWrapper;
  }

  public setValue(newValue: T) {
    this.variableWrapper = newValue;
  }
}
