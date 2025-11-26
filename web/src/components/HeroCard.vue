<template>
  <div class="hero min-h-screen">
    <div class="hero-content">
      <div class="card bg-base-100 shadow-md w-[90vw] max-w-xl">
        <div class="card-body h-[95vh] max-h-200 flex flex-col justify-between gap-5 lg:gap-10">
          <header>
            <slot name="header" />
          </header>
          <main class="flex-1 overflow-y-scroll">
            <slot />
          </main>
          <footer>
            <slot name="footer" />
          </footer>
        </div>
      </div>
    </div>
  </div>
</template>
