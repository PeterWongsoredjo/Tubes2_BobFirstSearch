export function MainTitle() {
  return (
    <div className="text-center mb-8 relative">
      <h1 className="font-cinzel text-5xl md:text-6xl lg:text-7xl text-amber-300  tracking-wide relative inline-block">
        <span className="text-6xl md:text-7xl lg:text-8xl">B</span>ob
        <span className="text-6xl md:text-7xl lg:text-8xl">F</span>irst
        <span className="text-6xl md:text-7xl lg:text-8xl">S</span>earch

        <span className="absolute -bottom-2 left-0 right-0 h-0.5 bg-gradient-to-r from-transparent via-amber-500/50 to-transparent"></span>
      </h1>

      <div className="mt-4 flex items-center justify-center gap-2">
        <p className="text-amber-200/80 max-w-2xl font-light tracking-wider">
          Mr. White we need to cook Mr. White!
        </p>
      </div>
      
    </div>
  )
}
