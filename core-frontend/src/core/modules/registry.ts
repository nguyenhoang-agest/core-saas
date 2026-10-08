import type { ModuleDefinition, ModuleMenuItem, SlotContribution } from './types';

export class ModuleRegistry {
  private modules: Map<string, ModuleDefinition> = new Map();

  /**
   * Registers a module into the registry.
   * Throws an error on duplicate module id or duplicate slot contribution id.
   */
  public register(module: ModuleDefinition): void {
    if (this.modules.has(module.id)) {
      throw new Error(`[ModuleRegistry] Duplicate module id detected: "${module.id}" is already registered.`);
    }

    if (module.slots) {
      const existingSlots = this.getAllSlotContributions();
      for (const slot of module.slots) {
        if (existingSlots.some((s) => s.id === slot.id)) {
          throw new Error(`[ModuleRegistry] Duplicate slot contribution id detected: "${slot.id}" already exists.`);
        }
      }
    }

    this.modules.set(module.id, module);
  }

  /**
   * Registers an array of modules.
   */
  public registerAll(modules: ModuleDefinition[]): void {
    modules.forEach((mod) => this.register(mod));
  }

  public get(id: string): ModuleDefinition | undefined {
    return this.modules.get(id);
  }

  public getAll(): ModuleDefinition[] {
    return Array.from(this.modules.values());
  }

  public getMenuItems(): ModuleMenuItem[] {
    const items: ModuleMenuItem[] = [];
    for (const mod of this.modules.values()) {
      if (mod.menuItems) {
        items.push(...mod.menuItems);
      }
    }
    return items.sort((a, b) => (a.order ?? 100) - (b.order ?? 100));
  }

  public getAllSlotContributions(): SlotContribution[] {
    const contributions: SlotContribution[] = [];
    for (const mod of this.modules.values()) {
      if (mod.slots) {
        contributions.push(...mod.slots);
      }
    }
    return contributions;
  }

  public getSlotContributions(slotName: string): SlotContribution[] {
    return this.getAllSlotContributions()
      .filter((contribution) => contribution.slotName === slotName)
      .sort((a, b) => (a.order ?? 100) - (b.order ?? 100));
  }

  public clear(): void {
    this.modules.clear();
  }
}

export const moduleRegistry = new ModuleRegistry();
